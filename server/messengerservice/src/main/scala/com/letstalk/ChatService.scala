package com.letstalk

import akka.actor.{ ActorRef, Props }
import com.letstalk.data_layer.DataManager

class ChatService
    extends ChatServer
    with SessionManagement
    with ChatManagement {

  implicit val storage: ActorRef = context.actorOf(DataManager.props(useMemory = true, useDatabase = false))

  override def preStart() = {
  }
}
