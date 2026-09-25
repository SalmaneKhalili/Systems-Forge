// Offline ambient soundscapes generated with Web Audio — no audio assets.
// A module-level singleton keeps the sound running across views (including
// the fullscreen deep-focus overlay).

export type AmbienceKind = "none" | "brown" | "rain" | "wind";

let ctx: AudioContext | null = null;
let master: GainNode | null = null;
let active: AudioScheduledSourceNode[] = [];
let lfoNodes: AudioNode[] = [];
let kind: AmbienceKind = "none";
let volume = 0.5;
let buffers: Record<string, AudioBuffer> = {};

function ensureCtx(): AudioContext | null {
  try {
    if (!ctx) ctx = new AudioContext();
    if (ctx.state === "suspended") void ctx.resume();
    return ctx;
  } catch {
    return null;
  }
}

function render(ctx: AudioContext, seconds: number, fn: (out: Float32Array) => void): AudioBuffer {
  const len = Math.floor(ctx.sampleRate * seconds);
  const buf = ctx.createBuffer(1, len, ctx.sampleRate);
  const data = buf.getChannelData(0);
  fn(data);
  return buf;
}

function whiteBuffer(ctx: AudioContext): AudioBuffer {
  if (buffers.white) return buffers.white;
  buffers.white = render(ctx, 4, (out) => {
    for (let i = 0; i < out.length; i++) out[i] = Math.random() * 2 - 1;
  });
  return buffers.white;
}

// Brown noise: integrated white noise (1/f²-ish), normalized.
function brownBuffer(ctx: AudioContext): AudioBuffer {
  if (buffers.brown) return buffers.brown;
  const white = whiteBuffer(ctx);
  buffers.brown = render(ctx, white.duration, (out) => {
    let last = 0;
    const src = white.getChannelData(0);
    for (let i = 0; i < out.length; i++) {
      last = (last + 0.02 * src[i]) / 1.02;
      out[i] = last * 3.5;
    }
  });
  return buffers.brown;
}

// Rain: steady noise bed + randomized decaying "droplet" taps baked into the
// buffer, plus a high shelf for brightness.
function rainBuffer(ctx: AudioContext): AudioBuffer {
  if (buffers.rain) return buffers.rain;
  const white = whiteBuffer(ctx);
  const src = white.getChannelData(0);
  buffers.rain = render(ctx, white.duration, (out) => {
    for (let i = 0; i < out.length; i++) {
      out[i] = src[i] * 0.35;
    }
    // sprinkle droplets: short exponential-decay blips at random offsets
    let i = 0;
    while (i < out.length) {
      i += Math.floor(200 + Math.random() * 4000);
      let j = 0;
      const amp = 0.12 + Math.random() * 0.2;
      while (j < 700 && i + j < out.length) {
        out[i + j] += (Math.random() * 2 - 1) * amp * Math.exp(-j / 90);
        j++;
      }
    }
  });
  return buffers.rain;
}

function stopAll(): void {
  for (const s of active) {
    try {
      s.stop();
    } catch {
      /* already stopped */
    }
  }
  for (const n of lfoNodes) {
    try {
      n.disconnect();
    } catch {
      /* noop */
    }
  }
  active = [];
  lfoNodes = [];
}

function connect(src: AudioScheduledSourceNode, ...chain: AudioNode[]): void {
  let node: AudioNode = src;
  for (const n of chain) {
    node.connect(n);
    node = n;
  }
  if (master) node.connect(master);
}

export function setAmbience(next: AmbienceKind, vol: number): void {
  const c = ensureCtx();
  if (!c || next === "none") {
    stopAll();
    kind = "none";
    volume = vol;
    return;
  }
  // (re)build master gain once
  if (!master) {
    master = c.createGain();
    master.connect(c.destination);
  }
  if (kind !== next) {
    stopAll();
    kind = next;
    if (next === "brown") {
      const b = brownBuffer(c);
      const src = c.createBufferSource();
      src.buffer = b;
      src.loop = true;
      const lp = c.createBiquadFilter();
      lp.type = "lowpass";
      lp.frequency.value = 800;
      connect(src, lp);
      active.push(src);
      src.start();
    } else if (next === "rain") {
      const b = rainBuffer(c);
      const src = c.createBufferSource();
      src.buffer = b;
      src.loop = true;
      const hp = c.createBiquadFilter();
      hp.type = "highpass";
      hp.frequency.value = 120;
      connect(src, hp);
      active.push(src);
      src.start();
    } else if (next === "wind") {
      const b = whiteBuffer(c);
      const src = c.createBufferSource();
      src.buffer = b;
      src.loop = true;
      const lp = c.createBiquadFilter();
      lp.type = "lowpass";
      lp.frequency.value = 900;
      const lfo = c.createOscillator();
      lfo.frequency.value = 0.15;
      const lfoGain = c.createGain();
      lfoGain.gain.value = 500;
      lfo.connect(lfoGain).connect(lp.frequency);
      connect(src, lp);
      active.push(src);
      lfoNodes.push(lfo, lfoGain);
      lfo.start();
      src.start();
    }
  }
  if (master && Number.isFinite(vol)) {
    master.gain.setTargetAtTime(Math.min(1, Math.max(0, vol)), c.currentTime, 0.1);
  }
  volume = vol;
}

export function ambienceKind(): AmbienceKind {
  return kind;
}

export function ambienceVolume(): number {
  return volume;
}

export const AMBIENCE_LABELS: { kind: AmbienceKind; label: string }[] = [
  { kind: "none", label: "Off" },
  { kind: "brown", label: "Brown noise" },
  { kind: "rain", label: "Rain" },
  { kind: "wind", label: "Wind" },
];