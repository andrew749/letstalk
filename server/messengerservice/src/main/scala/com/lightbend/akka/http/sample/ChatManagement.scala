package com.lightbend.akka.http.sample

import akka.actor.{Actor, ActorRef}
import com.lightbend.akka.http.sample.data_models.Message

import scala.collection.immutable.HashMap
import scala.collection.mutable

trait ChatManagement { this: Actor =>

  val sessions: mutable.HashMap[String, ActorRef]

  protected def chatManagement: Receive = {
    case msg @ Message(sender, recipient, payload) =>

      // send this message to the sessions if they exists
      getSession(recipient.id).foreach(_ ! msg)
      getSession(sender.id).foreach(_ ! msg)

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
      print(s"Message going to $userId")
      None
    }
  }
}
