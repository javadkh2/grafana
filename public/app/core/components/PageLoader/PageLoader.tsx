import React, { FC } from 'react';
import { translate } from '../../../locale/translator';

interface Props {
  pageName?: string;
}

const PageLoader: FC<Props> = ({ pageName = '' }) => {
  const loadingText = `${translate("Loading")} ${pageName}...`;
  return (
    <div className="page-loader-wrapper">
      <i className="page-loader-wrapper__spinner fa fa-spinner fa-spin" />
      <div className="page-loader-wrapper__text">{loadingText}</div>
    </div>
  );
};

export default PageLoader;
