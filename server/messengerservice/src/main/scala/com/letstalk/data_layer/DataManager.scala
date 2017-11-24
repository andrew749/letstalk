package com.letstalk.data_layer

import java.util.UUID

import akka.actor.{ Actor, ActorLogging, Props }
import com.letstalk.data_models.{ Message, UserModel }

import scala.collection.mutable

/**
 * This class receives messages with data that should be stored in
 * the appropriate datalayer type for it.
 *
 * Massive WIP
 */

class DataManager(useMemory: Boolean, useDatabase: Boolean) extends Actor with ActorLogging {

  implicit var dataLayers: mutable.Buffer[ChatStorage] = mutable.Buffer()

  if (useMemory) {
    val memoryChatStorage: MemoryChatStorage = new MemoryChatStorage()
    dataLayers.append(memoryChatStorage)
  }

  if (useDatabase) {
    val databaseChatStorage: DatabaseChatStorage = new DatabaseChatStorage()
    dataLayers.append(databaseChatStorage)
  }

  def receive: Receive = {
    case message: Message =>
      // store the message in all data layers
      dataLayers foreach { _ storeMessage message }

    case GetMessage(id) =>
      val results = dataLayers flatMap (_.retrieveMessage(id))

      sender() ! results.head

    case GetMessages(threadId) =>
      // This doesn't make much sense in the context of multiple datastores where you would need to
      // do some kind of consolidation so just using the first datastore for now.
      sender() ! dataLayers.head.retrieveMessages(threadId)

    case user: UserModel =>
      dataLayers foreach { _ storeUser user }

    case GetUser(id) =>
      val results = dataLayers flatMap (_.retrieveUser(id))
      sender() ! results.head
  }

}

object DataManager {
  def props(useMemory: Boolean, useDatabase: Boolean): Props = {
    Props(classOf[DataManager], useMemory, useDatabase)
  }
}

case class GetMessage(id: UUID)
case class GetUser(id: UUID)
case class GetMessages(threadId: UUID)
