package com.lightbend.akka.http.sample.routes

import akka.actor.{ActorRef, ActorSystem}
import akka.event.Logging
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.http.scaladsl.server.directives.MethodDirectives.{get, post}

/**
  * Routes which have to do with user actions
  */
trait UserRoutes {

  // we leave these abstract, since they will be provided by the App
  implicit def system: ActorSystem

  lazy val log = Logging(system, classOf[UserRoutes])

  // other dependencies that UserRoutes use
  def userRegistryActor: ActorRef

  lazy val userRoutes: Route =
    pathPrefix("users") {
      createUserRoute ~ getUserRoute
    }

  lazy val createUserRoute: Route = pathPrefix("create") {
    post {
      // TODO: send a message to the user registry actor to create a new user
      ???
    }
  }

  lazy val getUserRoute: Route = pathPrefix("get") {
    get {
      // TODO: get a users info
      ???
    }
  }
}
