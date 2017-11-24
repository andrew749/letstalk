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
    val uuid1 = UUID.randomUUID()
    val uuid2 = UUID.randomUUID()

    val user1 = NormalUser(uuid1, pinfo, cinfo)
    val user2 = NormalUser(uuid2, pinfo2, cinfo)
  }

  trait DataManagerTrait {
    val dataManager: ActorRef = system.actorOf(DataManager.props(useMemory = true, useDatabase = false))
  }

  implicit val timeout: Timeout = 10 seconds

  "A DataManager actor" must "store messages and return these messages" in new DataManagerTrait with TestUsers {

    val messagePayload = IncomingMessagePayload("Test Message", System.currentTimeMillis)
    val uuid = UUID.randomUUID()

    val message = Message(uuid, user1, user2, Option(messagePayload))

    dataManager ! message

    assert(Await.result(dataManager ? GetMessage(uuid), 10 seconds) === message)
  }

  "A DataManager actor" must "store users and return these user" in new DataManagerTrait with TestUsers {
    dataManager ! user1
    assert(Await.result(dataManager ? GetUser(user1.id), 10 seconds) === user1)
  }

  override protected def afterAll(): Unit = {
    TestKit.shutdownActorSystem(system)
  }
}
