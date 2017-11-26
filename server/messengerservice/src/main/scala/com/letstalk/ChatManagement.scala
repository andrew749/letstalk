package com.letstalk

import java.util.UUID

import akka.actor.{ Actor, ActorRef }
import akka.event.Logging

import com.letstalk.data_models.Message
import com.letstalk.data_layer.MessageEvent

case class WithAuth(token: UUID, messageEvent: MessageEvent)

trait ChatManagement { this: Actor =>

  // storage layer for the application
  val dataManager: ActorRef

  lazy val log = Logging.getLogger(context.system, this)

  protected def chatManagement: Receive = {
    case WithAuth(token, messageEvent) =>

      log.debug(s"Authenticating ${token}")
      // FIXME: Actually authenticate/authorize by calling authentication/authorization actors

      dataManager forward messageEvent
  }

}
