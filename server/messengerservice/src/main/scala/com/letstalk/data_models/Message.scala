package com.letstalk.data_models

import slick.model.Table

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

/**
 * An immutable message object to encapsulate data being sent
 * @param sender
 * @param recipient
 * @param payload
 */
final case class Message(sender: UserModel, recipient: UserModel, payload: Option[MessagePayload])

