package com.letstalk.routes

import java.util.UUID

import akka.actor.{ ActorRef, ActorSystem, Props }
import akka.event.Logging
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.pattern.ask
import akka.util.Timeout
import com.letstalk.UserRegistryActor.GetUser
import com.letstalk.data_models.{ IncomingMessagePayload, Message, Thread, UserModel }
import com.letstalk.data_layer.{ GetMessages, GetThreads, MemoryChatStorage, Messages, Threads }
import com.letstalk.{ ChatService, JsonSupport, WithAuth }

import scala.concurrent.Await
import scala.concurrent.duration._

case class MessageData(sender: UUID, thread: UUID, payload: String)

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

        path(Segment) { userId =>
          get {
            log.debug(s"Getting threads for ${userId}")
            val future = chatServerActor ? WithAuth(token, GetThreads(UUID.fromString(userId)))
            onSuccess(future) { case Threads(threads) => complete(threads) }
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
            complete(messageId)
          }
        }
      }
    }

  def getMessagesRoute(token: UUID): Route =
    pathPrefix("get") {
      path(Segment) { threadId =>
        get {
          val future = chatServerActor ? WithAuth(token, GetMessages(UUID.fromString(threadId)))
          onSuccess(future) { case Messages(msgs) => complete(msgs) }
        }
      }
    }
}
