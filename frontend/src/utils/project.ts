import { Project } from "../types";

export function projectName(project: Project) {
  let name = project.name;
  const postfix: string[] = [];
  if (project.rowStatus == "ARCHIVED") {
    postfix.push("Archived");
  }
  // if (project.tenantMode === "TENANT") {
  //   postfix.push("Tenant");
  // }

  if (postfix.length > 0) {
    name += ` (${postfix.join(", ")})`;
  }

  return name;
}
