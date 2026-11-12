<template>
	<div data-testid="hold-window-entry">
		<div v-if="showHeader" class="row d-none d-lg-flex mb-2">
			<div class="col-3">
				<label :for="formId('weekdays')">{{ $t("main.chargingPlan.weekdays") }}</label>
			</div>
			<div class="col-2">
				<label :for="formId('start')">{{ $t("main.chargingPlan.holdWindow.from") }}</label>
			</div>
			<div class="col-2">
				<label :for="formId('end')">{{ $t("main.chargingPlan.holdWindow.to") }}</label>
			</div>
			<div class="col-2">
				<label :for="formId('maxCurrent')">
					{{ $t("main.chargingPlan.holdWindow.maxCurrent") }}
				</label>
			</div>
			<div class="col-1">
				<label :for="formId('active')">{{ $t("main.chargingPlan.active") }}</label>
			</div>
		</div>
		<div class="row">
			<div class="col-5 d-lg-none col-form-label">
				<label :for="formId('weekdays')">{{ $t("main.chargingPlan.weekdays") }}</label>
			</div>
			<div class="col-7 col-lg-3 mb-2 mb-lg-0">
				<MultiSelect
					:id="formId('weekdays')"
					:value="selectedWeekdays"
					:options="dayOptions"
					:selectAllLabel="$t('main.chargingPlan.selectAll')"
					data-testid="hold-window-weekdays"
					@update:model-value="changeSelectedWeekdays"
				>
					{{ weekdaysLabel }}
				</MultiSelect>
			</div>
			<div class="col-5 d-lg-none col-form-label">
				<label :for="formId('start')">{{ $t("main.chargingPlan.holdWindow.from") }}</label>
			</div>
			<div class="col-7 col-lg-2 mb-2 mb-lg-0">
				<input
					:id="formId('start')"
					v-model="selectedStart"
					type="time"
					class="form-control mx-0 text-start"
					data-testid="hold-window-start"
					required
					@change="update()"
				/>
			</div>
			<div class="col-5 d-lg-none col-form-label">
				<label :for="formId('end')">{{ $t("main.chargingPlan.holdWindow.to") }}</label>
			</div>
			<div class="col-7 col-lg-2 mb-2 mb-lg-0">
				<input
					:id="formId('end')"
					v-model="selectedEnd"
					type="time"
					class="form-control mx-0 text-start"
					data-testid="hold-window-end"
					required
					@change="update()"
				/>
			</div>
			<div class="col-5 d-lg-none col-form-label">
				<label :for="formId('maxCurrent')">
					{{ $t("main.chargingPlan.holdWindow.maxCurrent") }}
				</label>
			</div>
			<div class="col-7 col-lg-2 mb-2 mb-lg-0">
				<select
					:id="formId('maxCurrent')"
					v-model="selectedMaxCurrent"
					class="form-select mx-0"
					data-testid="hold-window-max-current"
					@change="update()"
				>
					<option v-for="opt in currentOptions" :key="opt.value" :value="opt.value">
						{{ opt.name }}
					</option>
				</select>
			</div>
			<div class="col-5 d-lg-none col-form-label">
				<label :for="formId('active')">{{ $t("main.chargingPlan.active") }}</label>
			</div>
			<div class="col-3 col-lg-1 d-flex align-items-center">
				<div class="form-check form-switch">
					<input
						:id="formId('active')"
						v-model="selectedActive"
						class="form-check-input"
						type="checkbox"
						role="switch"
						data-testid="hold-window-active"
						tabindex="0"
						@change="update(true)"
					/>
				</div>
			</div>
			<div
				class="col-4 col-lg-2 d-flex align-items-center justify-content-end justify-content-lg-start"
			>
				<button
					v-if="showApply"
					type="button"
					class="btn btn-sm btn-outline-primary border-0 text-decoration-underline text-truncate"
					data-testid="hold-window-apply"
					tabindex="0"
					@click="update(true)"
				>
					{{ $t("main.chargingPlan.update") }}
				</button>
				<button
					v-else
					type="button"
					class="btn btn-sm btn-outline-secondary border-0"
					aria-label="Remove"
					data-testid="hold-window-remove"
					tabindex="0"
					@click="$emit('removed')"
				>
					<shopicon-regular-trash size="s" class="flex-shrink-0"></shopicon-regular-trash>
				</button>
			</div>
		</div>
	</div>
</template>

<script lang="ts">
import "@h2d2/shopicons/es/regular/trash";
import MultiSelect from "../Helper/MultiSelect.vue";
import formatter from "@/mixins/formatter";
import deepEqual from "@/utils/deepEqual";
import type { SelectOption } from "@/types/evcc";
import { defineComponent, type PropType } from "vue";

// hold power is meant to be low: offer common small currents, 0 = loadpoint maximum
const CURRENT_OPTIONS = [0, 6, 8, 10, 13, 16];

export default defineComponent({
	name: "ChargingPlanHoldWindowSettings",
	components: { MultiSelect },
	mixins: [formatter],
	props: {
		number: Number,
		weekdays: { type: Array as PropType<number[]>, default: () => [] },
		start: String,
		end: String,
		tz: String,
		minCurrent: { type: Number, default: 0 },
		maxCurrent: { type: Number, default: 0 },
		active: Boolean,
		showHeader: Boolean,
		formIdPrefix: String,
	},
	emits: ["updated", "removed"],
	data() {
		return {
			selectedWeekdays: this.weekdays,
			selectedStart: this.start,
			selectedEnd: this.end,
			selectedMaxCurrent: this.maxCurrent,
			selectedActive: this.active,
		};
	},
	computed: {
		dataChanged(): boolean {
			return (
				!deepEqual(this.weekdays, this.selectedWeekdays) ||
				this.start !== this.selectedStart ||
				this.end !== this.selectedEnd ||
				this.maxCurrent !== this.selectedMaxCurrent ||
				this.active !== this.selectedActive
			);
		},
		showApply(): boolean {
			return this.dataChanged && this.selectedActive;
		},
		weekdaysLabel(): string {
			return this.fmtWeekdaysRange(this.selectedWeekdays);
		},
		dayOptions(): SelectOption<number>[] {
			return this.getWeekdaysList("long");
		},
		currentOptions(): SelectOption<number>[] {
			return CURRENT_OPTIONS.map((value) => ({
				value,
				name: value === 0 ? "—" : `${value} A`,
			}));
		},
	},
	watch: {
		weekdays(newValue: number[], oldValue: number[]) {
			if (!deepEqual(newValue, oldValue)) {
				this.selectedWeekdays = newValue;
			}
		},
		start(newValue: string) {
			this.selectedStart = newValue;
		},
		end(newValue: string) {
			this.selectedEnd = newValue;
		},
		maxCurrent(newValue: number) {
			this.selectedMaxCurrent = newValue;
		},
		active(newValue: boolean) {
			this.selectedActive = newValue;
		},
	},
	methods: {
		changeSelectedWeekdays(weekdays: number[]): void {
			this.selectedWeekdays = weekdays;
			this.update();
		},
		formId(name: string): string {
			return `${this.formIdPrefix}-${this.number}-${name}`;
		},
		update(forceSave = false): void {
			const window = {
				weekdays: this.selectedWeekdays,
				start: this.selectedStart,
				end: this.selectedEnd,
				tz: this.tz,
				minCurrent: this.minCurrent,
				maxCurrent: this.selectedMaxCurrent,
				active: this.selectedActive,
			};

			if (forceSave || !this.selectedActive) {
				this.$emit("updated", window);
			}
		},
	},
});
</script>
