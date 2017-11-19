package com.letstalk

import akka.actor.{ Actor, ActorRef, Props }
import com.letstalk.sample.data_models.UserModel

import scala.collection.mutable

sealed trait UserEvent {
  val userId: String
}

case class LoginEvent(userId: String) extends UserEvent
case class SignupEvent(userId: String) extends UserEvent
case class LogoutEvent(userId: String) extends UserEvent

trait SessionManagement { this: Actor =>

  // storage layer for the application
  val storage: ActorRef

  // store all user sessions on this server
  val sessions = new mutable.HashMap[String, ActorRef]

  /**
   * Helper method to get a usermodel for a specific id
   * @param userId
   * @return
   */
  private def getUserForId(userId: String): UserModel = {
    ???
  }

  protected def sessionManagement: Receive = {
    case LoginEvent(userId) => {
      // get a user object
      val user = getUserForId(userId)

      val actor = context.actorOf(Props(new Session(user = user, storage = storage)))

      // add the session to the globally maintained list
      sessions += (userId -> actor)
    }

    case LogoutEvent(userId) => {
      // get the session corresponding to this userId
      val session = sessions(userId)

      // stop this actor
      context.stop(session)

      sessions -= userId
    }

    case SignupEvent(userId) => {
      ???
    }

  }

  protected def shutdownSessions = sessions.foreach { case (_, session) => context.stop(session) }

}
