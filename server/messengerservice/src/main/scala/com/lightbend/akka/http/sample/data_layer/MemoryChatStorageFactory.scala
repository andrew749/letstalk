package com.lightbend.akka.http.sample.data_layer

import akka.actor.{Actor, Props}

trait MemoryChatStorageFactory { this: Actor =>
  // storage to be used
  val storage = context.actorOf(Props[MemoryChatStorage])
}
