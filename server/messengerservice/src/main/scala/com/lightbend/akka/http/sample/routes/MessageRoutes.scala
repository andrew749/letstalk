package com.lightbend.akka.http.sample.routes

import akka.actor.{ActorRef, ActorSystem, Props}
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import com.lightbend.akka.http.sample.{ChatServer, JsonSupport}
import com.lightbend.akka.http.sample.data_models.IncomingMessagePayload
import akka.pattern.ask
import akka.util.Timeout

trait MessageRoutes extends JsonSupport  {

  implicit val system: ActorSystem
  implicit val timeout: Timeout

  // reference to the actor for the chat server
  lazy val chatServerActor: ActorRef = system.actorOf( Props[ChatServer] )

  lazy val messageRoute: Route =
    pathPrefix("messages") {
      sendMessageRoute
    }

  lazy val sendMessageRoute: Route =
    pathPrefix("send") {
      post {
        entity(as[IncomingMessagePayload]) { payload =>

          // send the message to the chat server
          val response = (chatServerActor ? payload).mapTo[String]

          complete(response)
        }
      }
    }

  lazy val getMessagesRoute: Route =
    pathPrefix("get") {
      path(Segment) { userId =>
       // TODO: get messages for this user
        ???
      }
    }

}
