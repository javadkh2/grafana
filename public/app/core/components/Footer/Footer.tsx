import React, { FC } from 'react';
import { Tooltip } from '@grafana/ui';
import { translate } from '../../../locale/translator';

interface Props {
  appName: string;
  buildVersion: string;
  buildCommit: string;
  newGrafanaVersionExists: boolean;
  newGrafanaVersion: string;
}

export const Footer: FC<Props> = React.memo(
  ({ appName, buildVersion, buildCommit, newGrafanaVersionExists, newGrafanaVersion }) => {
    return (
      <footer className="footer">
        <div className="text-center">
          <ul>
            <li>
              <a href="http://docs.grafana.org" target="_blank">
                <i className="fa fa-file-code-o" /> {translate("Docs")}
              </a>
            </li>
            <li>
              <a href="https://grafana.com/services/support" target="_blank">
                <i className="fa fa-support" /> {translate("Support Plans")}
              </a>
            </li>
            <li>
              <a href="https://community.grafana.com/" target="_blank">
                <i className="fa fa-comments-o" /> {translate("Community")}
              </a>
            </li>
            <li>
              <a href="https://grafana.com" target="_blank">
                {appName}
              </a>{' '}
              <span>
                v{buildVersion} (commit: {buildCommit})
              </span>
            </li>
            {newGrafanaVersionExists && (
              <li>
                <Tooltip placement="auto" content={newGrafanaVersion}>
                  <a href="https://grafana.com/get" target="_blank">
                    {translate("New version available!")}
                  </a>
                </Tooltip>
              </li>
            )}
          </ul>
        </div>
      </footer>
    );
  }
);

export default Footer;
