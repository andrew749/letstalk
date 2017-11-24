package com.letstalk.data_models

import java.util.UUID

import slick.jdbc.PostgresProfile.api._

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

/**
 * An immutable message object to encapsulate data being sent
 * @param id: unique identifier for the message
 * @param senderId: id of User object corresponding to who is sending the message.
 * @param threadId: id of Thread object that owns this message.
 * @param payload: Message data
 */
<<<<<<< HEAD
final case class Message(id: UUID, sender: UserModel, recipient: UserModel, payload: Option[MessagePayload])
=======
final case class Message(id: String, senderId: String, threadId: String, payload: Option[MessagePayload])
>>>>>>> WIP

/**
 * The table schema definition for a table.
 * @param tag
 */
<<<<<<< HEAD
class MessageTable(tag: Tag) extends Table[(UUID, UUID, UUID, UUID)](tag, "Messages") {
  def messageId = column[UUID]("ID", O.PrimaryKey)
  def senderId = column[UUID]("SENDER_ID")
  def recipientId = column[UUID]("RECIPIENT_ID")
  def payload = column[UUID]("PAYLOAD")

  // need to add this projection
  def * = (messageId, senderId, recipientId, payload)
=======
class MessageTable(tag: Tag) extends Table[(String, String, String, String)](tag, "Messages") {
  def messageId = column[String]("ID", O.PrimaryKey)
  def senderId = column[String]("SENDER_ID")
  def threadId = column[String]("THREAD_ID")
  def payload = column[String]("PAYLOAD")

  // need to add this projection
  def * = (messageId, senderId, threadId, payload)
>>>>>>> WIP

  val payloads = TableQuery[MessagePayloadTable]
  def payloadForeign = foreignKey("PAYLOAD_FK", payload, payloads)(_.id)
  val threadIds = TableQuery[ThreadTable]
  def threadIdForeign = foreignKey("THREAD_ID_FK", threadId, threadIds)(_.threadId)
  // TODO: User foreign key constraint when added
}
