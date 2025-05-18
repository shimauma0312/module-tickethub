import { mount } from "@vue/test-utils";
import { describe, expect, it } from "vitest";
import Index from '../pages/index.vue';

describe('Index Page', () => {
  it('ウェルカムメッセージが表示されること', () => {
    const wrapper = mount(Index);
    expect(wrapper.text()).toContain('TicketHub');
  });
});
