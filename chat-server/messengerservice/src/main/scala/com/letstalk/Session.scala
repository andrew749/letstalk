package com.letstalk

import akka.actor.{ Actor, ActorRef, Props }
import com.letstalk.data_models.{ Message, UserModel }

object Session {
}

class Session(user: UserModel, storage: ActorRef) extends Actor {

  private val startTime = System.currentTimeMillis
  private var log: List[Message] = Nil

  def receive = {

    case msg @ Message(id, sender, recipient, payload) =>
      // save this message into a in memory log
      log ::= msg

      // save this message into persistent storage
      storage forward msg
  }

}
