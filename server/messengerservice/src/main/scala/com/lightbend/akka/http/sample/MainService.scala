package com.lightbend.akka.http.sample

import akka.actor.{ActorRef, ActorSystem}
import akka.event.Logging
import akka.http.scaladsl.Http
import akka.http.scaladsl.Http.ServerBinding
import akka.http.scaladsl.server.Route
import akka.stream.ActorMaterializer
import akka.http.scaladsl.server.Directives._
import akka.util.Timeout
import com.lightbend.akka.http.sample.routes.{MessageRoutes, UserRoutes}

import scala.concurrent.duration._
import scala.concurrent.{ExecutionContext, Future}
import scala.io.StdIn

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
//#main-class
