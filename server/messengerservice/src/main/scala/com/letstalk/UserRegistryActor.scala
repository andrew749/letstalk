package com.letstalk

import akka.actor.{ Actor, ActorLogging, Props }
import com.letstalk.data_models.UserModel

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
      // add user to set
      // TODO: Save to storage with db logic
      log.debug(s"Creating user with id=${user.id}")
      users += user
      sender() ! ActionPerformed(s"User ${user.id} created.")

    case GetUser(id) =>
      log.debug(s"Checking for user with id=$id")
      sender() ! users.find(_.id == id)
  }
}
