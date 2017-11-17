package com.lightbend.akka.http.sample

import akka.actor.{Actor, ActorLogging, Props}
import akka.event.Logging
import com.lightbend.akka.http.sample.data_models.UserModel

//#user-case-classes

object UserRegistryActor {
  final case class ActionPerformed(description: String)
  final case class CreateUser(user: UserModel)
  final case class GetUser(id: String)

  def props: Props = Props[UserRegistryActor]
}

class UserRegistryActor extends Actor with ActorLogging {
  import UserRegistryActor._

  var users = Set.empty[UserModel]

  def receive: Receive = {

    case CreateUser(user) =>
      log.debug("Got user create message")
      // add user to set
      // TODO: Save to storage with db logic
      users += user
      sender() ! ActionPerformed(s"User ${user.id} created.")

    case GetUser(id) =>
      log.debug("Got user getuser message")
      sender() ! users.find(_.id == id)

    // TODO: Delete from storage with db logic

  }
}
