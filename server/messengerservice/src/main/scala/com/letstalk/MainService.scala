package com.letstalk

import akka.actor.{ActorRef, ActorSystem, Props}
import akka.event.Logging
import akka.http.scaladsl.Http
import akka.http.scaladsl.Http.ServerBinding
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.stream.ActorMaterializer
import akka.util.Timeout
import com.letstalk.UserRegistryActor.CreateUser
import com.letstalk.data_models.{ContactInfo, NormalUser, PersonalInfo}
import com.letstalk.routes.MessageRoutes
import com.letstalk.sample.routes.UserRoutes
import com.typesafe.config.ConfigFactory

import scala.concurrent.{ExecutionContext, Future}
import scala.concurrent.duration._
import scala.io.StdIn
import slick.jdbc.PostgresProfile.api._

//#main-class
object MainService extends App
    with UserRoutes
    with MessageRoutes {

  // set up ActorSystem and other dependencies here
  //#main-class
  val config = ConfigFactory.load()

  //#server-bootstrapping
  implicit val system: ActorSystem = ActorSystem("LetsTalk", config)
  implicit val materializer: ActorMaterializer = ActorMaterializer()
  //#server-bootstrapping

  private val log = Logging.getLogger(system, this)

  // database engine to use
  implicit val db = Database.forConfig("h2mem1")

  // Needed for the Future and its methods flatMap/onComplete in the end
  implicit val executionContext: ExecutionContext = system.dispatcher

  val timeout = Timeout(5 seconds)

  val userRegistryActor: ActorRef = system.actorOf(UserRegistryActor.props, "userRegistryActor")

  //#main-class
  // from the UserRoutes trait
  lazy val routes: Route = concat(
    userRoutes,
    messageRoute
  )

  //#main-class

  //#http-server
  val serverBindingFuture: Future[ServerBinding] = Http().bindAndHandle(routes, "localhost", 8080)

  log.info(s"Server online at http://localhost:8080/\nPress RETURN to stop...")

  StdIn.readLine()

  serverBindingFuture
    .flatMap(_.unbind())
    .onComplete { done =>
      done.failed.map { ex => log.error(ex, "Failed unbinding") }
      system.terminate()
    }
  //#http-server
  //#main-class
}
