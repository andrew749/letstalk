package com.letstalk.data_models

import java.util.UUID

import slick.jdbc.PostgresProfile.api._

import com.letstalk.data_layer.MessageEvent

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
final case class Message(id: UUID, senderId: UUID, threadId: UUID, payload: Option[MessagePayload])
  extends MessageEvent

/**
 * The table schema definition for a table.
 * @param tag
 */
class MessageTable(tag: Tag) extends Table[(UUID, UUID, UUID, UUID)](tag, "Messages") {
  def messageId = column[UUID]("ID", O.PrimaryKey)
  def senderId = column[UUID]("SENDER_ID")
  def threadId = column[UUID]("THREAD_ID")
  def payload = column[UUID]("PAYLOAD")

  // need to add this projection
  def * = (messageId, senderId, threadId, payload)

  val payloads = TableQuery[MessagePayloadTable]
  def payloadForeign = foreignKey("PAYLOAD_FK", payload, payloads)(_.id)
  val threadIds = TableQuery[ThreadTable]
  def threadIdForeign = foreignKey("THREAD_ID_FK", threadId, threadIds)(_.threadId)
  // TODO: User foreign key constraint when added
}
