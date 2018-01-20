package com.letstalk.data_layer

import java.util.UUID

import akka.actor.{ Actor, ActorLogging, Props }
import com.letstalk.data_models.{ Message, Thread, UserModel }

import scala.collection.mutable

trait MessageEvent

case class GetMessages(threadId: UUID) extends MessageEvent
case class GetThreads(userId: UUID) extends MessageEvent

case class Messages(values: List[Message])
case class Threads(values: List[Thread])

/**
 * This class receives messages with data that should be stored in
 * the appropriate datalayer type for it.
 *
 * Massive WIP
 */

class DataManager(storage: ChatStorage) extends Actor with ActorLogging {

  def receive: Receive = {
    case message: Message =>
      storage storeMessage message

    case GetMessages(threadId) =>
      val messages = storage retrieveMessages threadId
      sender() ! Messages(messages)

    case thread: Thread =>
      storage storeThread thread

    case GetThreads(userId) =>
      val threads = storage retrieveThreads userId
      sender() ! Threads(threads)
  }

}

object DataManager {
  def props(storage: ChatStorage): Props = {
    Props(classOf[DataManager], storage)
  }
}
