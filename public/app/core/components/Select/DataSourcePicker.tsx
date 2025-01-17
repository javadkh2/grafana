// Libraries
import React, { PureComponent } from 'react';

// Components
import { Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

// Types
import { DataSourceSelectItem } from '@grafana/ui';

import { translate } from '../../../locale/translator';

export interface Props {
  onChange: (ds: DataSourceSelectItem) => void;
  datasources: DataSourceSelectItem[];
  current: DataSourceSelectItem;
  onBlur?: () => void;
  autoFocus?: boolean;
  openMenuOnFocus?: boolean;
}

export class DataSourcePicker extends PureComponent<Props> {
  static defaultProps: Partial<Props> = {
    autoFocus: false,
    openMenuOnFocus: false,
  };

  searchInput: HTMLElement;

  constructor(props: Props) {
    super(props);
  }

  onChange = (item: SelectableValue<string>) => {
    const ds = this.props.datasources.find(ds => ds.name === item.value);
    this.props.onChange(ds);
  };

  render() {
    const { datasources, current, autoFocus, onBlur, openMenuOnFocus } = this.props;

    const options = datasources.map(ds => ({
      value: ds.name,
      label: ds.name,
      imgUrl: ds.meta.info.logos.small,
    }));

    const value = current && {
      label: current.name,
      value: current.name,
      imgUrl: current.meta.info.logos.small,
    };

    return (
      <div className="gf-form-inline">
        <Select
          className="ds-picker"
          isMulti={false}
          isClearable={false}
          backspaceRemovesValue={false}
          onChange={this.onChange}
          options={options}
          autoFocus={autoFocus}
          onBlur={onBlur}
          openMenuOnFocus={openMenuOnFocus}
          maxMenuHeight={500}
          placeholder={translate("Select datasource")}
          noOptionsMessage={() => translate("No datasources found")}
          value={value}
        />
      </div>
    );
  }
}

export default DataSourcePicker;
