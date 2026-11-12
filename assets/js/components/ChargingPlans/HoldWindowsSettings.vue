<template>
	<div data-testid="hold-windows">
		<div v-for="(window, index) in windows" :key="index" class="mb-4">
			<HoldWindowSettings
				:number="index + 1"
				:showHeader="index === 0"
				:formIdPrefix="formIdPrefix"
				v-bind="window"
				@updated="updateWindow(index, $event)"
				@removed="removeWindow(index)"
			/>
		</div>
		<div class="d-flex align-items-center pb-4">
			<button
				type="button"
				class="d-flex btn btn-sm btn-outline-secondary border-0 align-items-center gap-2 evcc-gray add-button"
				data-testid="hold-window-add"
				tabindex="0"
				@click="addWindow"
			>
				<shopicon-regular-plus size="s" class="flex-shrink-0"></shopicon-regular-plus>
				{{ $t("main.chargingPlan.holdWindow.add") }}
			</button>
		</div>
	</div>
</template>

<script lang="ts">
import "@h2d2/shopicons/es/regular/plus";
import HoldWindowSettings from "./HoldWindowSettings.vue";
import formatter from "@/mixins/formatter";
import { defineComponent, type PropType } from "vue";
import type { HoldWindow } from "@/types/evcc";

const DEFAULT_WEEKDAYS = [1, 2, 3, 4, 5];
const DEFAULT_START = "06:00";
const DEFAULT_END = "07:30";
const DEFAULT_MAX_CURRENT = 6;

export default defineComponent({
	name: "ChargingPlanHoldWindows",
	components: { HoldWindowSettings },
	mixins: [formatter],
	props: {
		id: [Number, String],
		windows: { type: Array as PropType<HoldWindow[]>, default: () => [] },
	},
	emits: ["updated"],
	computed: {
		formIdPrefix() {
			return `holdwindow-lp${this.id}`;
		},
	},
	methods: {
		addWindow(): void {
			const window: HoldWindow = {
				weekdays: DEFAULT_WEEKDAYS,
				start: DEFAULT_START,
				end: DEFAULT_END,
				tz: this.timezone(),
				minCurrent: 0,
				maxCurrent: DEFAULT_MAX_CURRENT,
				active: false,
			};
			this.$emit("updated", [...this.windows, window]);
		},
		updateWindow(index: number, window: HoldWindow): void {
			const windows = [...this.windows];
			windows.splice(index, 1, window);
			this.$emit("updated", windows);
		},
		removeWindow(index: number): void {
			const windows = [...this.windows];
			windows.splice(index, 1);
			this.$emit("updated", windows);
		},
	},
});
</script>

<style scoped>
.add-button {
	margin-left: -0.5rem;
}
</style>
