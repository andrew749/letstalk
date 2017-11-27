package com.letstalk.data_models

import java.util.UUID

import slick.jdbc.PostgresProfile.api._

import com.letstalk.data_layer.MessageEvent

/**
 * Immutable thread object that uniquely identifies a conversation.
 */
final case class Thread(id: UUID, userIds: List[UUID]) extends MessageEvent

/**
 * The table schema definition for the thread table
 * @param tag
 */
class ThreadTable(tag: Tag) extends Table[(UUID)](tag, "Threads") {
  def threadId = column[UUID]("ID", O.PrimaryKey)

  // need to add this projection
  def * = (threadId)
}
