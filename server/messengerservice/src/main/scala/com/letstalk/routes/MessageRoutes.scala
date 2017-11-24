package com.letstalk.routes

import java.util.UUID

import akka.actor.{ ActorRef, ActorSystem, Props }
import akka.event.Logging
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.pattern.ask
import akka.util.Timeout
import com.letstalk.UserRegistryActor.GetUser
import com.letstalk.data_models.{ IncomingMessagePayload, Message, UserModel }
import com.letstalk.data_layer.GetMessages
import com.letstalk.{ ChatService, JsonSupport }

import scala.concurrent.Await
import scala.concurrent.duration._

case class MessageData(from: String, to: String, payload: String)

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

          log.debug(s"Received message from ${data.from} to ${data.to}")

          implicit val atMost = 5 seconds
          // get futures for user data
          // FIXME: double trouble, first we block, second we cast
          val fromUser = Await.result(userRegistryActor ? GetUser(data.from), 3 seconds)
          val toThread = Await.result(chatServerActor ? GetThread(data.to), 3 seconds)
          (fromUser, toUser) match {
            case (Some(a: UserModel), Some(b: UserModel)) =>

              // get some id
              val tempId = getUUID()

              chatServerActor ! Message(tempId, a, b, Some(IncomingMessagePayload(data.payload, System.currentTimeMillis())))

              // return the generated id to the requested
              complete(tempId)

            case x =>
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

        val future = chatServerActor ? GetMessages(threadId)
        onSuccess(future) { messages: Seq[Message] => complete(messages.head) }
      }
    }

}
