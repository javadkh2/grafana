import React, { FC } from 'react';
import { translate } from '../../../locale/translator';

export const UserSignup: FC<{}> = () => {
  return (
    <div className="login-signup-box">
      <div className="login-signup-title p-r-1">{translate("New to Grafana?")}</div>
      <a href="signup" className="btn btn-medium btn-signup btn-p-x-2">
        {translate("Sign Up")}
      </a>
    </div>
  );
};
