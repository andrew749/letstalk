package com.letstalk.data_layer

import akka.actor.Actor
import com.letstalk.data_models.Message

import scala.collection.mutable

/**
  * This class receives messages with data that should be stored in
  * the appropriate datalayer type for it.
  *
  *
  * Massive WIP
  */
class DataManager extends Actor with MemoryChatStorage with DatabaseChatStorage {

  implicit var dataLayers: mutable.Buffer[ChatStorage] = mutable.Buffer()

  def receive: Receive = {
    case Message =>

  }

}


