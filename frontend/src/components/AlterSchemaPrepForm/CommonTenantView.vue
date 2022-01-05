<template>
  <div class="flex items-center py-1">
    <label for="project" class="textlabel mr-4">
      {{ $t("common.projects") }}
    </label>
    <div class="w-64">
      <ProjectSelect
        id="project"
        class="mt-1"
        name="project"
        :mode="ProjectMode.Tenant"
        :selected-id="local.projectId"
        @select-project-id="(id: number) => local.projectId = id"
      />
    </div>
  </div>
  <div v-if="local.projectId">
    <ProjectTenantView
      :state="state"
      :database-list="filteredDatabaseList"
      :environment-list="environmentList"
      :project="project"
      @dismiss="$emit('dismiss')"
    />
  </div>
</template>

<script lang="ts" setup>
import { defineProps, defineEmits, reactive, computed } from "vue";
import { useStore } from "vuex";
import { Database, Environment, Project, ProjectId } from "../../types";
import ProjectSelect, { Mode as ProjectMode } from "../ProjectSelect.vue";
import ProjectTenantView, {
  State as ProjectTenantState,
} from "./ProjectTenantView.vue";

export type State = ProjectTenantState;

const props = defineProps<{
  state: State;
  databaseList: Database[];
  environmentList: Environment[];
  project?: Project;
}>();

defineEmits<{
  (event: "dismiss"): void;
}>();

const store = useStore();

const local = reactive({
  projectId: undefined as ProjectId | undefined,
});

const project = computed(() => {
  return store.getters["project/projectById"](local.projectId) as Project;
});

const filteredDatabaseList = computed(() => {
  if (!local.projectId) return [];
  return props.databaseList.filter((db) => db.project.id === local.projectId);
});
</script>
