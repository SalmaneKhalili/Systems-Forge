/// <reference types="svelte" />

declare module "markdown-it" {
  const MarkdownIt: any;
  export default MarkdownIt;
}
declare module "highlight.js" {
  const hljs: any;
  export default hljs;
}

interface WailsApp {
  Curriculum(): Promise<Curriculum>;
  GetExercise(id: string): Promise<ExerciseDetail>;
  ListFiles(id: string): Promise<FileInfo[]>;
  ReadFile(id: string, rel: string): Promise<string>;
  WriteFile(id: string, rel: string, content: string): Promise<void>;
  RunCheck(id: string): Promise<CheckResult>;
  Progress(): Promise<Progress>;
  QuizDeck(filter: string): Promise<Card[]>;
  AnswerCard(key: string, rating: number): Promise<Card>;
  StartFocus(kind: string, exID: string): Promise<StudySession>;
  EndFocus(id: number, minutes: number): Promise<StudySession>;
  Stats(): Promise<Stats>;
  GetSettings(): Promise<Record<string, string>>;
  SaveSettings(kv: Record<string, string>): Promise<void>;
}

interface Window {
  go?: {
    main?: {
      App: WailsApp;
    };
  };
}