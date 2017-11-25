package com.letstalk

import java.util.UUID

import akka.actor.{ Actor, ActorRef }
import akka.event.Logging
import com.letstalk.data_models.Message

import scala.collection.mutable

trait ChatManagement { this: Actor =>

  val sessions: mutable.HashMap[UUID, ActorRef]

  lazy val log = Logging.getLogger(context.system, this)

  protected def chatManagement: Receive = {
    case msg @ Message(id, senderId, threadId, payload) =>
      log.debug("Got Message!")

      // send this message to the sessions if they exists
      // NOTE: mixing threads and users, not too sure about this
      getSession(threadId) foreach { _ ! msg }
      getSession(senderId) foreach { _ ! msg }

  }

  /**
   * Return any exising session for this user
   * @param userId
   * @return Option[ActorRef] The (possibly non-existent) session for this user
   */
  private def getSession(userId: UUID): Option[ActorRef] = {
    if (sessions.contains(userId)) {
      Some(sessions(userId))
    } else {
      log.debug(s"Message going to $userId")
      None
    }
  }
}
