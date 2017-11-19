package com.letstalk

trait ChatServer extends Actor {
  val storage: ActorRef

  // compose partial functions
  def receive: Receive = sessionManagement orElse chatManagement

  // abstract methods defined elsewhere
  protected def sessionManagement: Receive
  protected def chatManagement: Receive
  protected def shutdownSessions: Unit


  override def postStop(): Unit = {
    // delete all sessions
    shutdownSessions

    context.stop(storage)
  }
}
