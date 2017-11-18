package com.lightbend.akka.http.sample

import akka.actor.{Actor, ActorRef}
import akka.event.Logging
import com.lightbend.akka.http.sample.data_models.Message

import scala.collection.mutable

trait ChatManagement { this: Actor =>

  val sessions: mutable.HashMap[String, ActorRef]

  lazy val log = Logging.getLogger(context.system, this)

  protected def chatManagement: Receive = {
    case msg @ Message(sender, recipient, payload) =>
      log.debug("Got Message!")

      // send this message to the sessions if they exists
      getSession(recipient.id) foreach {_ ! msg}
      getSession(sender.id) foreach {_ ! msg}

  }

  /**
   * Return any exising session for this user
   * @param userId
   * @return Option[ActorRef] The (possibly non-existent) session for this user
   */
  private def getSession(userId: String): Option[ActorRef] = {
    if (sessions.contains(userId)) {
      Some(sessions(userId))
    } else {
      log.debug(s"Message going to $userId")
      None
    }
  }
}
