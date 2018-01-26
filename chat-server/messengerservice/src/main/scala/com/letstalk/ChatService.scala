package com.letstalk

import akka.actor.{ ActorRef, Props }
import com.letstalk.data_layer.{ ChatStorage, DataManager }

class ChatService(storage: ChatStorage)
    extends ChatServer
    with SessionManagement
    with ChatManagement {

  implicit val dataManager: ActorRef = context.actorOf(DataManager.props(storage))
}

object ChatService {
  def props(storage: ChatStorage): Props = {
    Props(classOf[ChatService], storage)
  }
}
