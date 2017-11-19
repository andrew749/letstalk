package com.letstalk

import com.letstalk.sample.data_layer.MemoryChatStorageFactory

class ChatService
    extends ChatServer
    with SessionManagement
    with ChatManagement
    with MemoryChatStorageFactory {

  override def preStart() = {
  }
}
