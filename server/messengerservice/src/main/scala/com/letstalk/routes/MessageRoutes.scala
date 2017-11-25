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
import com.letstalk.data_layer.{ GetMessages, GetThread }
import com.letstalk.{ ChatService, JsonSupport }

import scala.concurrent.Await
import scala.concurrent.duration._

case class MessageData(sender: UUID, thread: UUID, payload: String)

case class Messages(values: Seq[Message])

trait MessageRoutes extends JsonSupport {

  implicit val system: ActorSystem
  implicit val timeout: Timeout

  def getUUID(): UUID = {
    java.util.UUID.randomUUID
  }

  // reference to the actor for the chat server
  lazy val chatServerActor: ActorRef = system.actorOf(Props[ChatService])

  // TODO: make this a public service
  val userRegistryActor: ActorRef

  private lazy val log = Logging.getLogger(system, this)

  lazy val messageRoute: Route =
    pathPrefix("messages") {
      sendMessageRoute
    }

  lazy val sendMessageRoute: Route =
    pathPrefix("send") {

      // messages sent via post
      post {

        log.debug("post")
        entity(as[MessageData]) { data: MessageData =>

          log.debug(s"Received message from ${data.sender} to ${data.thread}")

          implicit val atMost = 5 seconds
          // get futures for user data
          // FIXME: double trouble, first we block, second we cast
          val sender = Await.result(userRegistryActor ? GetUser(data.sender), 3 seconds)
          val thread = Await.result(chatServerActor ? GetThread(data.thread), 3 seconds)
          (sender, thread) match {
            case (Some(_: UserModel), Some(_: Thread)) =>

              // get some id
              val tempId = getUUID()

              chatServerActor ! Message(tempId, data.sender, data.thread,
                Some(IncomingMessagePayload(data.payload, System.currentTimeMillis())))

              // return the generated id to the requested
              complete(tempId)

            case x =>
              // FIXME: Proper error handling, perhaps the thread or user don't exist
              log.debug(s"Got an unknown type ${x}")
          }

          complete("OK")
        }
      }

    }

  lazy val getMessagesRoute: Route =
    pathPrefix("get") {
      path(Segment) { threadId =>
        log.debug("Got message get route")

        val future = chatServerActor ? GetMessages(UUID.fromString(threadId))
        onSuccess(future) {
          case Messages(msgs) => complete(msgs)
          case x => complete("Unknown Type")
        }
      }
    }

}
