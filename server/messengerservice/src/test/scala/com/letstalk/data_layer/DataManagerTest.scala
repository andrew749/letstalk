package com.letstalk.data_layer

import java.util.UUID

import akka.actor.{ ActorRef, ActorSystem, Props }
import akka.pattern.ask
import akka.testkit.{ ImplicitSender, TestKit }
import akka.util.Timeout
import com.letstalk.data_models._
import org.scalatest.{ BeforeAndAfterAll, FlatSpecLike, MustMatchers }

import scala.concurrent.Await
import scala.concurrent.duration._

class DataManagerTest() extends TestKit(ActorSystem("DataManagerTest"))
    with FlatSpecLike
    with MustMatchers
    with ImplicitSender
    with BeforeAndAfterAll {

  trait TestUsers {
    val pinfo = PersonalInfo("acod")
    val pinfo2 = PersonalInfo("Andrew 2")
    val cinfo = ContactInfo("test@gmail.com", "5555555555")
    val user1UUID = UUID.randomUUID()
    val user2UUID = UUID.randomUUID()
    val thread1UUID = UUID.randomUUID()

    val user1 = NormalUser(user1UUID, pinfo, cinfo)
    val user2 = NormalUser(user2UUID, pinfo, cinfo)
    val thread1 = Thread(thread1UUID, user1UUID :: user2UUID :: Nil)
  }

  trait DataManagerTrait {
    val dataManager: ActorRef = system.actorOf(DataManager.props(new MemoryChatStorage))
  }

  implicit val timeout: Timeout = 10 seconds

  "A DataManager actor" must "store messages and return these messages" in new DataManagerTrait with TestUsers {

    val messagePayload = IncomingMessagePayload("Test Message", System.currentTimeMillis)
    val uuid = UUID.randomUUID()

    val message = Message(uuid, user1UUID, thread1UUID, Option(messagePayload))

    dataManager ! message

    assert(Await.result(dataManager ? GetMessages(thread1UUID), 10 seconds) === Messages(Seq(message)))
  }

  "A DataManager actor" must "store threads and returns these threads" in new DataManagerTrait with TestUsers {
    dataManager ! thread1
    assert(Await.result(dataManager ? GetThreads(user1.id), 10 seconds) === Threads(Seq(thread1)))
    assert(Await.result(dataManager ? GetThreads(user2.id), 10 seconds) === Threads(Seq(thread1)))
  }

  override protected def afterAll(): Unit = {
    TestKit.shutdownActorSystem(system)
  }
}
