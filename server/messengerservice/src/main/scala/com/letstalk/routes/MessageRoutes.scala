package com.letstalk.routes

import java.util.UUID

import akka.actor.{ActorRef, ActorSystem}
import akka.event.Logging
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.pattern.ask
import akka.util.Timeout
import com.letstalk.data_layer._
import com.letstalk.data_models.{CreateThread, IncomingMessagePayload, Message, Thread}
import com.letstalk.{ChatService, JsonSupport, WithAuth}

import scala.concurrent.duration._

case class MessageData(sender: UUID, thread: UUID, payload: String)
case class NewThread()

case class SendMessageResponse(messageId: UUID)

trait MessageRoutes extends JsonSupport {

  implicit val system: ActorSystem
  implicit val timeout: Timeout

  def getUUID(): UUID = {
    UUID.randomUUID
  }

  // reference to the actor for the chat server
  lazy val chatServerActor: ActorRef = system.actorOf(ChatService.props(new MemoryChatStorage))

  // TODO: make this a public service
  val userRegistryActor: ActorRef

  private lazy val log = Logging.getLogger(system, this)

  lazy val messageRoute: Route =
    // /messages
    pathPrefix("messages") {
      cookie("sid") { sid =>
        val token = UUID.fromString(sid.value)
        sendMessageRoute(token) ~
          getMessagesRoute(token)
      }
    }

  lazy val threadRoute: Route =
    pathPrefix("threads") {
      // FIXME: Place session cookie stuff outside of thread/messageRoute
      cookie("sid") { sid =>
        val token = UUID.fromString(sid.value)

        // /threads/<userId>
        path(Segment) { userId =>
          get {
            log.debug(s"Getting threads for ${userId}")
            val threadsFuture = chatServerActor ? WithAuth(token, GetThreads(UUID.fromString(userId)))
            onSuccess(threadsFuture) { case Threads(threads) => complete(threads) }
          }
        }

        // /threads/create
        pathPrefix("create") {
          post {
            entity(as[CreateThread]) { newThread: CreateThread =>

              val threadId = getUUID()
              log.debug(s"Creating new thread with Id=${threadId}")

              // FIXME: Check if a users exist and error if they do not
              chatServerActor ! WithAuth(token, Thread(threadId, newThread.userIds))
              complete(threadId)
            }
          }
        }

      }
    }

  def sendMessageRoute(token: UUID): Route =
    pathPrefix("send") {

      post {

        log.debug("post")

        entity(as[MessageData]) { data: MessageData =>

          cookie("sid") { sid =>
            // FIXME: probably don't need to send sender since that can be found from session token
            log.debug(s"Received message from ${data.sender} to ${data.thread}")
            implicit val atMost = 5 seconds

            // FIXME: Not entirely sure that this belongs here, probably should have a separate
            // CreateMessage event
            val messageId = getUUID()
            val message = Message(messageId, data.sender, data.thread,
              Some(IncomingMessagePayload(data.payload, System.currentTimeMillis())))

            chatServerActor ! WithAuth(token, message)

            // return the generated id to the requested
            complete(SendMessageResponse(messageId))
          }
        }
      }
    }

  def getMessagesRoute(token: UUID): Route =
    pathPrefix("get") {
      // /messages/get/<threadId>
      path(Segment) { threadId =>
        get {
          val messagesFuture = chatServerActor ? WithAuth(token, GetMessages(UUID.fromString(threadId)))
          onSuccess(messagesFuture) { case Messages(msgs) => complete(msgs) }
        }
      }
    }
}
