package com.lightbend.akka.http.sample.routes

import akka.actor.{ActorRef, ActorSystem, Props}
import akka.event.Logging
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.pattern.ask
import akka.util.Timeout
import com.lightbend.akka.http.sample.UserRegistryActor.GetUser
import com.lightbend.akka.http.sample.data_models.{IncomingMessagePayload, Message, UserModel}
import com.lightbend.akka.http.sample.{ChatService, JsonSupport}

import scala.concurrent.ExecutionContext.Implicits.global
import scala.concurrent.Future

case class MessageData(from: String, to: String, payload: String)

trait MessageRoutes extends JsonSupport {

  implicit val system: ActorSystem
  implicit val timeout: Timeout

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

      post {
        log.debug("Got message send route")
        entity(as[IncomingMessagePayload]) { payload =>

          // send the message to the chat server
          val response = (chatServerActor ? payload).mapTo[String]

          complete(response)
        }
      }

      get {
        entity(as[MessageData]) { data: MessageData =>
          log.debug("Got message")

          val fromUser: Future[UserModel] = (userRegistryActor ? GetUser(data.from)).mapTo[UserModel]
          val toUser: Future[UserModel] = (userRegistryActor ? GetUser(data.to)).mapTo[UserModel]

          fromUser foreach {
            from =>
              {
                toUser foreach {
                  to =>
                    chatServerActor ! Message(from, to, Some(IncomingMessagePayload(data.payload, System.currentTimeMillis())))
                }
              }
          }

          complete("OK")
        }
      }

    }

  lazy val getMessagesRoute: Route =
    pathPrefix("get") {
      path(Segment) { userId =>
        log.debug("Got message get route")
        // TODO: get messages for this user
        ???
      }
    }

}
