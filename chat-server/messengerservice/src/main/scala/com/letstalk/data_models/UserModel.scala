package com.letstalk.data_models

import java.util.UUID

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

/**
 * Base trait to represent a user
 */
trait UserModel {
  def id: UUID
  def info: PersonalInfo
  def contactInfo: ContactInfo
}

/**
 * Specific user of the platform
 * @param id identifier of this user
 * @param info personal information to identify a user
 * @param contactInfo
 */
final case class NormalUser(id: UUID, info: PersonalInfo, contactInfo: ContactInfo) extends UserModel {}

/**
 * Administrator of the platform
 *
 * Need a differentiation since an administrator account can perform special actions and have access to special interface.
 * @param id identifier of this user
 * @param info personal information to identify a user
 * @param contactInfo
 */
final case class AdministratorUser(id: UUID, info: PersonalInfo, contactInfo: ContactInfo) extends UserModel {}