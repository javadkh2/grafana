import React from 'react';
import { translate } from '../../../locale/translator';

export default () => {
  return (
    <div className="">
      <h5>{translate("What are Permissions?")}</h5>
      <p>
        {translate("An Access Control List (ACL) model is used to limit access to Dashboard Folders.")}
        {translate("A user or a Team can be assigned permissions for a folder or for a single dashboard.")}
      </p>
    </div>
  );
};
