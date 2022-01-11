<template>
  <main class="flex-1 relative overflow-y-auto">
    <!-- Highlight Panel -->
    <div class="px-4 pt-2">
      <!-- Summary -->
      <div class="flex items-center">
        <div>
          <div class="flex items-center">
            <h1
              class="pt-2 pb-2.5 text-md font-bold leading-6 text-main truncate"
            >
              {{ database.name }}
            </h1>
          </div>
        </div>
      </div>
      <dl class="flex flex-col space-y-2 md:space-y-0 md:flex-row md:flex-wrap">
        <dt class="sr-only">{{ $t("common.environment") }}</dt>
        <dd class="flex items-center text-sm md:mr-4">
          <span class="textlabel"
            >{{ $t("common.environment") }}&nbsp;-&nbsp;</span
          >
          <router-link
            :to="`/environment/${environmentSlug(
              database.instance.environment
            )}`"
            class="normal-link"
            >{{ environmentName(database.instance.environment) }}</router-link
          >
        </dd>
        <dt class="sr-only">{{ $t("common.instance") }}</dt>
        <dd class="flex items-center text-sm md:mr-4">
          <InstanceEngineIcon :instance="database.instance" />
          <span class="ml-1 textlabel"
            >{{ $t("common.instance") }}&nbsp;-&nbsp;</span
          >
          <router-link
            :to="`/instance/${instanceSlug(database.instance)}`"
            class="normal-link"
            >{{ instanceName(database.instance) }}</router-link
          >
        </dd>
        <template v-if="database.sourceBackup">
          <dt class="sr-only">{{ $t("db.parent") }}</dt>
          <dd class="flex items-center text-sm md:mr-4 tooltip-wrapper">
            <span class="textlabel">
              {{ $t("database.restored-from") }}
            </span>
            <router-link
              :to="`/db/${database.sourceBackup.databaseId}`"
              class="normal-link"
            >
              <!-- Do not display the name of the backup's database because that requires a fetch  -->
              <span class="tooltip">
                {{
                  $t(
                    "database.database-name-is-restored-from-another-database-backup",
                    [database.name]
                  )
                }}
              </span>
              {{ $t("database.database-backup") }}
            </router-link>
          </dd>
        </template>
        <dd class="flex items-center text-sm md:mr-4">
          <label class="textlabel">
            {{ $t("common.labels") }}&nbsp;-&nbsp;
          </label>
          <DatabaseLabels :labels="database.labels" :editable="false" />
        </dd>
        <dd class="flex items-center text-sm md:mr-4">
          <span class="textlabel">{{ $t("sql-editor.self") }}</span>
          <button class="ml-1 btn-icon" @click.prevent="gotoSqlEditor">
            <heroicons-outline:terminal class="w-4 h-4" />
          </button>
        </dd>
      </dl>
    </div>
  </main>
</template>

<script lang="ts" setup>
import { defineProps } from "vue";
import { useRouter } from "vue-router";
import { Database } from "../../types";
import { connectionSlug } from "../../utils";
import { DatabaseLabels } from "../DatabaseLabels/";

const props = defineProps<{
  database: Database;
}>();

const router = useRouter();

const gotoSqlEditor = () => {
  router.push({
    name: "sql-editor.detail",
    params: {
      connectionSlug: connectionSlug(props.database),
    },
  });
};
</script>
