package com.lightbend.akka.http.sample

import com.lightbend.akka.http.sample.data_layer.MemoryChatStorageFactory

class ChatService
    extends ChatServer
    with SessionManagement
    with ChatManagement
    with MemoryChatStorageFactory {

  override def preStart() = {
  }
}
