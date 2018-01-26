package com.letstalk.data_models

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

sealed trait Info

/**
 * Encapsulate personal information of a user.
 * @param name name of a user
 */
final case class PersonalInfo(name: String) extends Info {}

/**
 * Encapsulate contact information for a user.
 * @param email email address of a user
 */
final case class ContactInfo(email: String, phone: String) extends Info {}