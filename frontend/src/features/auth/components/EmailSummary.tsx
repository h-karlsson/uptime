import { Button } from '@/components/ui/button'

/** Shows the email chosen in the first step, with a way back to change it. */
export function EmailSummary({ email, onChange }: { email: string; onChange: () => void }) {
  return (
    <span className="inline-flex items-center gap-1">
      <span className="font-medium text-foreground">{email}</span>
      <Button type="button" variant="link" size="xs" onClick={onChange}>
        Change
      </Button>
    </span>
  )
}
