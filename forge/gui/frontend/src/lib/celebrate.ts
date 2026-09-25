// Zero-asset celebrations: a canvas confetti burst + a synthesized Web Audio
// chime. Both are generated locally, so everything stays offline.

let audioCtx: AudioContext | null = null;

function ctx(): AudioContext | null {
  try {
    if (!audioCtx) audioCtx = new AudioContext();
    if (audioCtx.state === "suspended") void audioCtx.resume();
    return audioCtx;
  } catch {
    return null;
  }
}

// ding plays a short, pleasant two-note chime (success feedback).
export function ding(): void {
  const c = ctx();
  if (!c) return;
  const t = c.currentTime;
  const notes = [523.25, 783.99]; // C5 → G5
  notes.forEach((freq, i) => {
    const osc = c.createOscillator();
    const gain = c.createGain();
    osc.type = "sine";
    osc.frequency.value = freq;
    const start = t + i * 0.09;
    gain.gain.setValueAtTime(0.0001, start);
    gain.gain.exponentialRampToValueAtTime(0.18, start + 0.02);
    gain.gain.exponentialRampToValueAtTime(0.0001, start + 0.5);
    osc.connect(gain).connect(c.destination);
    osc.start(start);
    osc.stop(start + 0.55);
  });
}

// confetti fires a burst from the top-center of the viewport. Call in a
// browser context (after the view is mounted).
export function confetti(amount = 120): void {
  const canvas = document.createElement("canvas");
  canvas.style.cssText =
    "position:fixed;inset:0;pointer-events:none;z-index:9999;";
  document.body.appendChild(canvas);
  const dpr = window.devicePixelRatio || 1;
  canvas.width = innerWidth * dpr;
  canvas.height = innerHeight * dpr;
  const g = canvas.getContext("2d");
  if (!g) {
    canvas.remove();
    return;
  }
  g.scale(dpr, dpr);

  const colors = ["#f43f5e", "#f59e0b", "#10b981", "#3b82f6", "#a855f7", "#facc15"];
  const parts = Array.from({ length: amount }, (_, i) => {
    const angle = Math.random() * Math.PI * 2;
    const speed = 4 + Math.random() * 7;
    return {
      x: innerWidth / 2,
      y: innerHeight * 0.3,
      vx: Math.cos(angle) * speed,
      vy: Math.sin(angle) * speed - 5,
      w: 6 + Math.random() * 6,
      h: 8 + Math.random() * 6,
      rot: Math.random() * Math.PI,
      vr: (Math.random() - 0.5) * 0.3,
      color: colors[i % colors.length],
    };
  });

  let t = 0;
  const gravity = 0.22;
  const drag = 0.985;
  const frame = () => {
    t++;
    g.clearRect(0, 0, innerWidth, innerHeight);
    let alive = false;
    for (const p of parts) {
      if (p.y > innerHeight + 40) continue;
      alive = true;
      p.vy += gravity;
      p.vx *= drag;
      p.vy *= drag;
      p.x += p.vx;
      p.y += p.vy;
      p.rot += p.vr;
      g.save();
      g.translate(p.x, p.y);
      g.rotate(p.rot);
      g.fillStyle = p.color;
      g.fillRect(-p.w / 2, -p.h / 2, p.w, p.h);
      g.restore();
    }
    if (alive && t < 240) {
      requestAnimationFrame(frame);
    } else {
      canvas.remove();
    }
  };
  requestAnimationFrame(frame);
}