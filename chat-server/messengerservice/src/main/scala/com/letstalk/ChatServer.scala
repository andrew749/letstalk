package com.letstalk

import akka.actor.{ Actor, ActorRef }

trait ChatServer extends Actor {
  val dataManager: ActorRef

  // compose partial functions
  def receive: Receive = sessionManagement orElse chatManagement

  // abstract methods defined elsewhere
  protected def sessionManagement: Receive
  protected def chatManagement: Receive

  override def postStop(): Unit = {
    // delete all sessions
    context.stop(dataManager)
  }

  override def preStart(): Unit = {}
}
