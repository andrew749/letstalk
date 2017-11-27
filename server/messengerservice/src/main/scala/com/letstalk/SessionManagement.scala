package com.letstalk

import java.util.UUID

import akka.actor.{ Actor, ActorRef, Props }
import com.letstalk.data_models.UserModel

import scala.collection.mutable

sealed trait UserEvent {
  val userId: UUID
}

case class LoginEvent(userId: UUID) extends UserEvent
case class SignupEvent(userId: UUID) extends UserEvent
case class LogoutEvent(userId: UUID) extends UserEvent

trait SessionManagement { this: Actor =>

  // storage layer for the application
  val dataManager: ActorRef

  protected def sessionManagement: Receive = {
    case LoginEvent(userId) => {
      ???
    }

    case LogoutEvent(userId) => {
      ???
    }

    case SignupEvent(userId) => {
      ???
    }

  }

}
