import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput, Switch, FieldSet } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { GizmoSQLDataSourceOptions, GizmoSQLSecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<GizmoSQLDataSourceOptions, GizmoSQLSecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData, secureJsonFields, secureJsonData } = options;

  const onHostChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        host: event.target.value,
      },
    });
  };

  const onPortChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        port: parseInt(event.target.value, 10) || 31337,
      },
    });
  };

  const onUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        username: event.target.value,
      },
    });
  };

  const onTLSChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        useTLS: event.currentTarget.checked,
      },
    });
  };

  const onSkipTLSVerifyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        skipTLSVerify: event.currentTarget.checked,
      },
    });
  };

  const onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...secureJsonData,
        password: event.target.value,
      },
    });
  };

  const onResetPassword = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...secureJsonFields,
        password: false,
      },
      secureJsonData: {
        ...secureJsonData,
        password: '',
      },
    });
  };

  const onTokenChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...secureJsonData,
        token: event.target.value,
      },
    });
  };

  const onResetToken = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...secureJsonFields,
        token: false,
      },
      secureJsonData: {
        ...secureJsonData,
        token: '',
      },
    });
  };

  return (
    <>
      <FieldSet label="Connection">
        <InlineField label="Host" labelWidth={14} tooltip="GizmoSQL server hostname or IP address">
          <Input
            width={40}
            value={jsonData.host || ''}
            onChange={onHostChange}
            placeholder="localhost"
          />
        </InlineField>
        <InlineField label="Port" labelWidth={14} tooltip="GizmoSQL Flight SQL port (default: 31337)">
          <Input
            width={20}
            type="number"
            value={jsonData.port || 31337}
            onChange={onPortChange}
            placeholder="31337"
          />
        </InlineField>
        <InlineField label="Use TLS" labelWidth={14} tooltip="Enable TLS/SSL encryption">
          <Switch value={jsonData.useTLS || false} onChange={onTLSChange} />
        </InlineField>
        {jsonData.useTLS && (
          <InlineField label="Skip TLS Verify" labelWidth={14} tooltip="Skip TLS certificate verification (not recommended for production)">
            <Switch value={jsonData.skipTLSVerify || false} onChange={onSkipTLSVerifyChange} />
          </InlineField>
        )}
      </FieldSet>

      <FieldSet label="Authentication">
        <InlineField label="Username" labelWidth={14} tooltip="Username for authentication (optional)">
          <Input
            width={40}
            value={jsonData.username || ''}
            onChange={onUsernameChange}
            placeholder="username"
          />
        </InlineField>
        <InlineField label="Password" labelWidth={14} tooltip="Password for authentication (optional)">
          <SecretInput
            width={40}
            isConfigured={secureJsonFields?.password ?? false}
            value={secureJsonData?.password || ''}
            onChange={onPasswordChange}
            onReset={onResetPassword}
            placeholder="password"
          />
        </InlineField>
        <InlineField label="Token" labelWidth={14} tooltip="Bearer token for authentication (optional, alternative to password)">
          <SecretInput
            width={40}
            isConfigured={secureJsonFields?.token ?? false}
            value={secureJsonData?.token || ''}
            onChange={onTokenChange}
            onReset={onResetToken}
            placeholder="bearer token"
          />
        </InlineField>
      </FieldSet>
    </>
  );
}

