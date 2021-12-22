import { Project } from ".";
import { DeploymentConfigId } from "./id";
import { LabelKey, LabelValue } from "./label";
import { Principal } from "./principal";

export type DeploymentConfig = {
  id: DeploymentConfigId;

  // Standard fields
  creator: Principal;
  createdTs: number;
  updater: Principal;
  updatedTs: number;

  // Related fields
  project: Project;

  // Domain specific fields
  schedule: DeploymentSchedule;
};

export type DeploymentSchedule = {
  deployments: Deployment[];
};

export type Deployment = {
  spec: DeploymentSpec;
};

export type DeploymentSpec = {
  selector: LabelSelector;
};

export type LabelSelector = {
  matchExpressions: LabelSelectorRequirement[];
};

export type LabelSelectorRequirement = {
  key: LabelKey;
  operator: OperatorType;
  values: LabelValue[];
};

export type OperatorType = "In" | "Exists";
