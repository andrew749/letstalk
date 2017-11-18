package com.lightbend.akka.http.sample.routes

import akka.actor.{ActorRef, ActorSystem}
import akka.event.Logging
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.http.scaladsl.server.directives.MethodDirectives.{get, post}
import akka.pattern.ask
import akka.util.Timeout
import com.lightbend.akka.http.sample.JsonSupport
import com.lightbend.akka.http.sample.UserRegistryActor.GetUser
import com.lightbend.akka.http.sample.data_models.{AdministratorUser, NormalUser, UserModel}

import scala.concurrent.{Await, Future}
import scala.concurrent.duration._

/**
 * Routes which have to do with user actions
 */
trait UserRoutes extends JsonSupport {

  // we leave these abstract, since they will be provided by the App
  implicit def system: ActorSystem

  private lazy val log = Logging.getLogger(system, this)

  // other dependencies that UserRoutes use
  def userRegistryActor: ActorRef

  lazy val userRoutes: Route =
    pathPrefix("users") {
      createUserRoute ~ getUserRoute
    }

  lazy val createUserRoute: Route = pathPrefix("create") {
    post {
      log.debug("Got user create route")
      // TODO: send a message to the user registry actor to create a new user
      ???
    }
  }

  lazy val getUserRoute: Route = pathPrefix("get") {
    get {
      parameters('test.as[String]) { id =>
        // TODO: perform auth on this request
        log.debug("Got user get route")
        // TODO: get a users info
        val duration = 10 seconds
        implicit val timeout = Timeout(10 seconds)
        val future = userRegistryActor ? GetUser(id)
        onSuccess(future) {
          case Some(x) =>
            // FIXME: this is always going to default case.
            x match {
              case user: NormalUser =>
                complete("normal user")
              case user: AdministratorUser =>
                complete("admin user")
              case _ =>
                complete("what")
            }
          case None => complete("Not found")
        }
      }
    }
  }
}
