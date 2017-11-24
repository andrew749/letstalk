package com.letstalk.data_models

import slick.jdbc.PostgresProfile.api._

/**
 * Immutable thread object that uniquely identifies a conversation.
 */
final case class Thread(id: String)

/**
 * The table schema definition for the thread table
 * @param tag
 */
class ThreadTable(tag: Tag) extends Table[(String)](tag, "Threads") {
  def threadId = column[String]("ID", O.PrimaryKey)

  // need to add this projection
  def * = (threadId)
}
