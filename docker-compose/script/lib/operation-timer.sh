#!/usr/bin/env bash

start_operation_timer() {
  OPERATION_TIMER_LABEL="$1"
  OPERATION_TIMER_START=$SECONDS
  trap print_operation_elapsed_time EXIT
}

print_operation_elapsed_time() {
  local exit_code=$?
  local elapsed=$((SECONDS - OPERATION_TIMER_START))
  local hours=$((elapsed / 3600))
  local minutes=$(((elapsed % 3600) / 60))
  local seconds=$((elapsed % 60))
  local duration

  if [ "$hours" -gt 0 ]; then
    duration="${hours}小时${minutes}分${seconds}秒"
  elif [ "$minutes" -gt 0 ]; then
    duration="${minutes}分${seconds}秒"
  else
    duration="${seconds}秒"
  fi

  if [ "$exit_code" -eq 0 ]; then
    echo "${OPERATION_TIMER_LABEL}完成，总耗时: $duration"
  else
    echo "${OPERATION_TIMER_LABEL}失败，总耗时: $duration" >&2
  fi
}
