package com.lightbend.akka.http.sample

class ChatService
  extends ChatServer
    with SessionManagement
    with ChatManagement
    with MemoryChatStorageFactory {

  override def preStart() = {
  }
}
